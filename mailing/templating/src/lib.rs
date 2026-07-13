use libc::c_char;
use mrml::prelude::render::RenderOptions;
use std::ffi::{CStr, CString};

fn ptr_to_str(ptr: *const c_char) -> &'static str {
    let c_str = unsafe { CStr::from_ptr(ptr) };

    c_str.to_str().expect("This should not be empty")
}

#[unsafe(no_mangle)]
pub extern "C" fn render_mrml_template(tmpl: *const c_char) -> *const c_char {
    let msg = mrml::parse(ptr_to_str(tmpl))
        .expect("The template should be parseable")
        .element
        .render(&RenderOptions::default())
        .expect("The template is invalid");

    CString::new(msg.as_str())
        .expect("The message should have been rendered")
        .into_raw()
}
