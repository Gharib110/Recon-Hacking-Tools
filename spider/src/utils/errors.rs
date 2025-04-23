use thiserror::Error;

#[derive(Error, Debug, Clone)]
pub enum Error {
    #[error("Internal")]
    Internal(String),
    #[error("Spider is not valid: {0}")]
    InvalidSpider(String),
    #[error("Request: {0}")]
    Request(String),
    #[error("WebDriver: {0}")]
    WebDriver(String),
}

impl std::convert::From<reqwest::Error> for Error {
    fn from(err: reqwest::Error) -> Self {
        Error::Request(err.to_string())
    }
}

impl std::convert::From<fantoccini::error::CmdError> for Error {
    fn from(err: fantoccini::error::CmdError) -> Self {
        Error::WebDriver(err.to_string())
    }
}

impl std::convert::From<fantoccini::error::NewSessionError> for Error {
    fn from(err: fantoccini::error::NewSessionError) -> Self {
        Error::WebDriver(err.to_string())
    }
}