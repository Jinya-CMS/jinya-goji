use env_logger::Env;
use sea_orm::sea_query::prelude::rust_decimal::prelude::ToPrimitive;
use std::ffi::{CStr, CString, c_char, c_int, c_ushort};
use std::str::Utf8Error;
use std::sync::OnceLock;
use tokio::runtime::Runtime;

pub mod dbal;
pub mod http;
pub mod models;
pub mod theming;

static RUNTIME: OnceLock<Runtime> = OnceLock::new();

fn runtime() -> &'static Runtime {
    RUNTIME.get_or_init(|| Runtime::new().expect("failed to create Tokio runtime"))
}

fn ptr_to_str(ptr: *const c_char) -> Result<String, Utf8Error> {
    let c_str = unsafe { CStr::from_ptr(ptr) };

    c_str.to_str().map(|res| res.to_string())
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn handle_request(
    api_request: bool,
    url: *mut c_char,
    output: *mut *mut c_char,
) -> c_int {
    let env = Env::new().filter("LOG_LEVEL");

    let _ = env_logger::try_init_from_env(env);

    let url = if let Ok(url) = ptr_to_str(url) {
        url
    } else {
        return 1;
    };

    unsafe {
        *output = std::ptr::null_mut();
    }

    runtime().block_on(async {
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
pub unsafe extern "C" fn handle_error(
    api_request: bool,
    url: *mut c_char,
    status_code: c_ushort,
    output: *mut *mut c_char,
) -> c_int {
    let env = Env::new().filter("LOG_LEVEL");

    let _ = env_logger::try_init_from_env(env);

    let url = if let Ok(url) = ptr_to_str(url) {
        url
    } else {
        return 1;
    };

    unsafe {
        *output = std::ptr::null_mut();
    }

    runtime().block_on(async {
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
pub unsafe extern "C" fn free_engine_string(ptr: *mut c_char) {
    if ptr.is_null() {
        return;
    }

    unsafe {
        drop(CString::from_raw(ptr));
    }
}
