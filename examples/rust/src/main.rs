use std::{
    error::Error,
    ffi::{CStr, CString},
    fs::read_to_string,
    path::PathBuf,
};

use clap::Parser;

use crate::soratun::Send;

#[allow(non_snake_case)]
#[allow(dead_code)]
#[allow(non_camel_case_types)]
#[allow(unused_qualifications)]
mod soratun;

#[derive(Parser)]
struct Args {
    /// Path to the soratun configuration file.
    #[clap(short, long, default_value = "arc.json")]
    config: PathBuf,

    /// HTTP method. Only GET or POST (case insensitive) is supported.
    #[clap(short, long, default_value = "POST")]
    method: String,

    /// HTTP path.
    #[clap(short, long, default_value = "/")]
    path: String,

    /// HTTP body.
    #[clap()]
    body: String,
}

fn main() -> Result<(), Box<dyn Error>> {
    let Args {
        config,
        method,
        path,
        body,
    } = Args::parse();

    let config = CString::new(read_to_string(config)?)?.into_raw();
    let method = CString::new(method)?.into_raw();
    let path = CString::new(path)?.into_raw();
    let body = CString::new(body)?.into_raw();

    let result = unsafe {
        let r = Send(config, method, path, body);
        CStr::from_ptr(r).to_str()?
    };
    println!("{result}");

    Ok(())
}
