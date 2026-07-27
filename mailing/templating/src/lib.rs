pub mod render;

use crate::render::{render_template, MrmlError};
use std::ffi::{c_char, c_int, CStr, CString};
use std::str::Utf8Error;
use env_logger::Env;

fn ptr_to_str(ptr: *const c_char) -> Result<String, Utf8Error> {
    let c_str = unsafe { CStr::from_ptr(ptr) };

    c_str.to_str().map(|res| res.to_string())
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn render_mrml_template(tmpl: *const c_char, output: *mut *mut c_char) -> c_int {
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
    let result = render_template(&tmpl);
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
        Err(MrmlError::ParseError) => 2,
        Err(MrmlError::RenderError) => 3,
    }
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn free_mrml_string(ptr: *mut c_char) {
    if ptr.is_null() {
        return;
    }

    unsafe {
        drop(CString::from_raw(ptr));
    }
}
