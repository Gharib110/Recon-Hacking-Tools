use clap::{Command, Arg};
use std::{sync::Arc, time::Duration};
use crate::crawler::Crawler;
use utils::errors::*;

mod crawler;
mod utils;


#[tokio::main]
async fn main() -> Result<(), anyhow::Error> {
    let cli = Command::new(clap::crate_name!())
        .version(clap::crate_version!())
        .about(clap::crate_description!())
        .subcommand(Command::new("spiders").about("List all spiders"))
        .subcommand(
            Command::new("run").about("Run a spider").arg(
                Arg::new("spider")
                    .short('s')
                    .long("spider")
                    .help("The spider to run")
                    .takes_value(true)
                    .required(true),
            ),
        )
        .arg_required_else_help(true)
        .get_matches();

    if let Some(_) = cli.subcommand_matches("spiders") {
        let spider_names = vec!["cvedetails", "github", "quotes"];
        for name in spider_names {
            println!("{}", name);
        }
    } else if let Some(matches) = cli.subcommand_matches("run") {
        // we can safely unwrap as the argument is required
        let spider_name = matches.value_of("spider").unwrap();
        let crawler = Crawler::new(Duration::from_millis(200), 2, 500);

        match spider_name {
            "cvedetails" => {
                let spider = Arc::new(utils::cve_details::CveDetailsSpider::new());
                crawler.run(spider).await;
            }
            "github" => {
                let spider = Arc::new(utils::github::GitHubSpider::new());
                crawler.run(spider).await;
            }
            "quotes" => {
                let spider = utils::quotes::QuotesSpider::new().await?;
                let spider = Arc::new(spider);
                crawler.run(spider).await;
            }
            _ => return Err(Error::InvalidSpider(spider_name.to_string()).into()),
        };
    }

    Ok(())
}