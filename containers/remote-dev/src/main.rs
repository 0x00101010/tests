use std::thread;
use std::time::Duration;

fn main() {
    loop {
        println!("Hello");
        thread::sleep(Duration::from_secs(5));
    }
}
