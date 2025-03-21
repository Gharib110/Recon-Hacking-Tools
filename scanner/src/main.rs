use std::env;
use std::error::Error;
use std::time::Duration;
use rayon::prelude::IntoParallelIterator;
use reqwest::{redirect, Client};
use crate::errors::Error::ScannerUsage;
use crate::utils::{scan_ports, subdomain_enumeration, Subdomain};

mod errors;
mod utils;

fn main() -> Result<(), Box<dyn Error>> {
    let args: Vec<String> = env::args().collect();

    if args.len() != 2 {
        return Err(ScannerUsage.into());
    }

    let target = args[1].as_str();

    let http_timeout = Duration::from_secs(5);
    let http_client = Client::builder()
        .redirect(redirect::Policy::limited(4))
        .timeout(http_timeout)
        .build()?;


    let pool = rayon::ThreadPoolBuilder::new()
        .num_threads(25)
        .build()
        .unwrap();


    pool.install(|| {
        let scan_result: Vec<Subdomain> = subdomain_enumeration(&http_client, target)
            .unwrap()
            .into_par_iter()
            .map(scan_ports)
            .collect();

        for subdomain in scan_result {
            println!("{}:", &subdomain.domain);
            for port in &subdomain.open_ports {
                println!("    {}", port.port);
            }

            println!();
        }
    });

    Ok(())
}
