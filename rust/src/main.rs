use std::{
    error::Error,
    ffi::{
        c_char,
        CStr,
        CString,
    },
    fs::File,
    io::Read,
};

use clap::Parser;
use crate::soratun::Send;

#[allow(non_snake_case)]
#[allow(dead_code)]
#[allow(non_camel_case_types)]
#[allow(unused_qualifications)]
mod soratun;

#[derive(Parser, Debug)]
#[clap()]
struct Args {
    /// Path to the Soracom Arc config file.
    #[clap(short, long, default_value = "arc.json")]
    config: String,

    /// HTTP method.
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
    let (config, method, path, body) = into_raw(Args::parse())?;

    let result = unsafe {
        let r = Send(config, method, path, body);
        CStr::from_ptr(r).to_str()?
    };
    println!("{result}");

    Ok(())
}

fn into_raw(args: Args) -> Result<(*mut c_char, *mut c_char, *mut c_char, *mut c_char), Box<dyn Error>> {
    let config = read_config(&args.config)?.into_raw();
    let method = CString::new(args.method)?.into_raw();
    let path = CString::new(args.path)?.into_raw();
    let body = CString::new(args.body)?.into_raw();
    Ok((config, method, path, body))
}

fn read_config(path: &str) -> Result<CString, Box<dyn Error>> {
    let mut file = File::open(path)?;
    let mut contents = String::new();
    file.read_to_string(&mut contents)?;
    Ok(CString::new(contents)?)
}
