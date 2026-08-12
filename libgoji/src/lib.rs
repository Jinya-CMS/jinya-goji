use engine::{http, preview};
use env_logger::Env;
use num_traits::ToPrimitive;
use std::ffi::{CStr, CString, c_char, c_int, c_ushort};
use std::str::Utf8Error;
use std::sync::OnceLock;
use tokio::runtime::Runtime;

fn ptr_to_str(ptr: *const c_char) -> Result<String, Utf8Error> {
    let c_str = unsafe { CStr::from_ptr(ptr) };

    c_str.to_str().map(|res| res.to_string())
}

static RUNTIME: OnceLock<Runtime> = OnceLock::new();

fn run_async<F>(future: F) -> c_int
where
    F: Future<Output = c_int>,
{
    let env = Env::new().filter("LOG_LEVEL");

    let _ = env_logger::try_init_from_env(env);

    RUNTIME
        .get_or_init(|| Runtime::new().expect("failed to create Tokio runtime"))
        .block_on(future)
}

#[unsafe(no_mangle)]
pub extern "C" fn compiler_compile_theme(theme_id: c_int) -> c_int {
    run_async(async move {
        match compiler::compiler::compile_theme(theme_id).await {
            Ok(()) => 0,
            Err(error) => {
                log::error!("Failed to compile theme {theme_id}: {error}");
                1
            }
        }
    })
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn mailing_render_mrml_template(
    tmpl: *const c_char,
    output: *mut *mut c_char,
) -> c_int {
    let env = Env::new().filter("LOG_LEVEL");

    let _ = env_logger::try_init_from_env(env);

    unsafe {
        *output = std::ptr::null_mut();
    }

    let tmpl = if let Ok(tmpl) = ptr_to_str(tmpl) {
        tmpl
    } else {
        return 1;
    };
    let result = mailing::render::render_template(&tmpl);
    match result {
        Ok(rendered) => {
            let result = CString::new(rendered.as_str())
                .expect("This should not happen")
                .into_raw();
            unsafe {
                *output = result;
            }
            0
        }
        Err(mailing::render::MrmlError::ParseError) => 2,
        Err(mailing::render::MrmlError::RenderError) => 3,
    }
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn engine_handle_request(
    api_request: bool,
    url: *const c_char,
    output: *mut *mut c_char,
) -> c_int {
    let url = if let Ok(url) = ptr_to_str(url) {
        url
    } else {
        return 1;
    };

    unsafe {
        *output = std::ptr::null_mut();
    }

    run_async(async move {
        let (status_code, body) = http::handle_request(&url, api_request).await;
        let result = CString::new(body.as_str())
            .expect("This should not happen")
            .into_raw();
        unsafe {
            *output = result;
        }

        status_code.as_u16().to_i32().unwrap_or(500)
    })
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn engine_handle_error(
    api_request: bool,
    url: *const c_char,
    status_code: c_ushort,
    output: *mut *mut c_char,
) -> c_int {
    let url = if let Ok(url) = ptr_to_str(url) {
        url
    } else {
        return 1;
    };

    unsafe {
        *output = std::ptr::null_mut();
    }

    run_async(async move {
        let (status_code, body) = http::handle_error(&url, api_request, status_code).await;
        let result = CString::new(body.as_str())
            .expect("This should not happen")
            .into_raw();
        unsafe {
            *output = result;
        }

        status_code.as_u16().to_i32().unwrap_or(500)
    })
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn engine_preview_classic_page(
    theme_id: c_int,
    classic_page_id: c_int,
    classic_page_data: *const c_char,
    output: *mut *mut c_char,
) -> c_int {
    let classic_page_data = if let Ok(classic_page_data) = ptr_to_str(classic_page_data) {
        classic_page_data
    } else {
        return 1;
    };

    unsafe {
        *output = std::ptr::null_mut();
    }

    run_async(async move {
        match preview::classic_page::classic_page_preview(
            theme_id,
            classic_page_id,
            classic_page_data,
        )
        .await
        {
            Ok(body) => {
                let result = CString::new(body.as_str())
                    .expect("This should not happen")
                    .into_raw();
                unsafe {
                    *output = result;
                }

                0
            }
            Err(err) => {
                log::error!("Failed to render classic page preview {err}");
                -1
            }
        }
    })
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn engine_preview_modern_page(
    theme_id: c_int,
    modern_page_id: c_int,
    modern_page_sections_data: *const c_char,
    output: *mut *mut c_char,
) -> c_int {
    let modern_page_sections_data =
        if let Ok(modern_page_sections_data) = ptr_to_str(modern_page_sections_data) {
            modern_page_sections_data
        } else {
            return 1;
        };

    unsafe {
        *output = std::ptr::null_mut();
    }

    run_async(async move {
        match preview::modern_page::modern_page_preview(
            theme_id,
            modern_page_id,
            modern_page_sections_data,
        )
        .await
        {
            Ok(body) => {
                let result = CString::new(body.as_str())
                    .expect("This should not happen")
                    .into_raw();
                unsafe {
                    *output = result;
                }

                0
            }
            Err(err) => {
                log::error!("Failed to render modern page preview {err}");
                -1
            }
        }
    })
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn engine_preview_blog_post(
    theme_id: c_int,
    blog_post_id: c_int,
    blog_post_sections_data: *const c_char,
    output: *mut *mut c_char,
) -> c_int {
    let blog_post_sections_data =
        if let Ok(blog_post_sections_data) = ptr_to_str(blog_post_sections_data) {
            blog_post_sections_data
        } else {
            return 1;
        };

    unsafe {
        *output = std::ptr::null_mut();
    }

    run_async(async move {
        match preview::blog_post::blog_post_preview(theme_id, blog_post_id, blog_post_sections_data)
            .await
        {
            Ok(body) => {
                let result = CString::new(body.as_str())
                    .expect("This should not happen")
                    .into_raw();
                unsafe {
                    *output = result;
                }

                0
            }
            Err(err) => {
                log::error!("Failed to render blog post preview {err}");
                -1
            }
        }
    })
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn free_string(ptr: *mut c_char) {
    if ptr.is_null() {
        return;
    }

    unsafe {
        drop(CString::from_raw(ptr));
    }
}
