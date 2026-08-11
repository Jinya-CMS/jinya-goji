use wasmtime::component::bindgen;

pub mod engine;
pub mod frontend_engine;
pub mod preview_engine;
pub mod theme;

bindgen!({
    path: "wit/goji.wit",
    imports: {
        default: async
    },
    exports: {
        default: async
    }
});
