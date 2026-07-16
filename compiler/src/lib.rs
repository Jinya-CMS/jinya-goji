pub mod compiler;
pub mod dbal;

use env_logger::Env;
use std::ffi::c_int;

#[unsafe(no_mangle)]
pub extern "C" fn compile_theme(theme_id: c_int) -> c_int {
    let env = Env::new().filter("LOG_LEVEL");

    let _ = env_logger::try_init_from_env(env);

    let runtime = match tokio::runtime::Runtime::new() {
        Ok(runtime) => runtime,
        Err(error) => {
            eprintln!("Failed to create Tokio runtime: {error}");
            return 2;
        }
    };

    runtime.block_on(async {
        match compiler::compile_theme(theme_id).await {
            Ok(()) => 0,
            Err(error) => {
                log::error!("Failed to compile theme {theme_id}: {error}");
                1
            }
        }
    })
}
