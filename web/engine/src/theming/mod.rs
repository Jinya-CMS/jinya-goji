use wasmtime::component::bindgen;

pub mod theme;
pub mod engine;

bindgen!({
    path: "wit/goji.wit",
    imports: {
        default: async
    },
    exports: {
        default: async
    }
});