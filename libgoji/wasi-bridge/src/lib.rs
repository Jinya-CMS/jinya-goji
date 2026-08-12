use wasmtime::component::bindgen;

bindgen!({
    path: "../../theming/wit/goji.wit",
    imports: {
        default: async
    },
    exports: {
        default: async
    }
});
