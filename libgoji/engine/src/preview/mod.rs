use std::path::Path;
use wasmtime::{Engine, Store};
use wasmtime::component::{Component, HasSelf, Linker, ResourceTable};
use wasmtime_wasi::WasiCtx;
use wasi_bridge::GojiTheme;
use crate::theming::preview_engine::{GojiPreviewEngineStore, PreviewEntity};
use crate::theming::theme::Theme;

pub mod classic_page;
pub mod modern_page;
pub mod blog_post;

async fn get_goji_theme(
    theme: &Theme,
    preview_entity: PreviewEntity,
    route: String,
) -> anyhow::Result<(GojiTheme, Store<GojiPreviewEngineStore>)> {
    let theme_path = Path::new(&theme.db_theme.path)
        .join("theme.wasm")
        .to_string_lossy()
        .to_string();

    let engine = Engine::default();
    let mut linker = Linker::new(&engine);
    wasmtime_wasi::p2::add_to_linker_async(&mut linker).map_err(|err| {
        log::error!(
            "Failed to append wasip2 for theme {}: {err}",
            theme.db_theme.name
        );
        err
    })?;

    let wasi_ctx = WasiCtx::builder().inherit_stdio().inherit_args().build();

    let mut store = Store::new(
        &engine,
        GojiPreviewEngineStore {
            theme: theme.clone(),
            uri: route.clone(),
            path: route,
            wasi_ctx,
            resource_table: ResourceTable::new(),
            preview_entity,
        },
    );

    let component = Component::from_file(&engine, &theme_path).map_err(|err| {
        log::error!(
            "Failed to build component for theme {}: {err:#?}",
            theme.db_theme.name
        );
        err
    })?;

    GojiTheme::add_to_linker::<_, HasSelf<GojiPreviewEngineStore>>(&mut linker, |state| state)
        .map_err(|err| {
            log::error!(
                "Failed to add goji theme {} to linker: {err}",
                theme.db_theme.name
            );
            err
        })?;

    Ok((
        GojiTheme::instantiate_async(&mut store, &component, &linker)
            .await
            .map_err(|err| {
                log::error!(
                    "Failed to instantiate goji theme {}: {err}",
                    theme.db_theme.name
                );
                err
            })?,
        store,
    ))
}
