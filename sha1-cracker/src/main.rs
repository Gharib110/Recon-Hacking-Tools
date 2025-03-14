use std::error::Error;
use std::fs::File;
use std::io::{BufRead, BufReader};
use sha1::Digest;

const SHA1_HEX_STRING_LENGTH: usize = 40;

fn main() -> Result<(), Box<dyn Error>> {
   let args: Vec<String> = std::env::args().collect();
    if args.len() != 3 {
        std::println!("Usage:");
        std::println!("cargo run -- <wordlist-path> <sha1-hash>"); //println! prevent format-string-vuln
        return Ok(());
    }

    let hash = args[2].trim();
    if hash.len() != SHA1_HEX_STRING_LENGTH {
        return Err("sha1 hash is not valid".into());
    }

    let wordlist = File::open(&args[1])?;
    let wordlist_reader = BufReader::new(&wordlist);

    for line in wordlist_reader.lines() {
        let line = line?;
        let common_password = line.trim();
        if hash == &hex::encode(sha1::Sha1::digest(common_password.as_bytes())) {
            println!("Password found: {}", &common_password);
            return Ok(());
        }
    }
    println!("password not found in wordlist :(");

    // Run: cargo run -- wordlist.txt 1790afb503c6c1e5d2f42af0d86adbdd75d20561
    // It is the hash of Helllo
    return Ok(());
}
