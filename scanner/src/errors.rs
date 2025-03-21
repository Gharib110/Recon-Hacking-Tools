use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error {
    #[error("Usage: Scanner <alirezagharib.ir>")]
    ScannerUsage,
}