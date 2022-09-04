mod database;
use database::Database;

fn main() {
    let mut database_path = String::from(".ok");
    if cfg!(windows) {
        database_path = format!("C:\\Users\\{}\\AppData\\Roaming\\ok", whoami::username())
    } else if cfg!(unix) {
        database_path = format!("/home/{}/.config/ok", whoami::username())
    }
    let db = match Database::open(database_path) {
        Ok(db) => db,
        Err(error) => {
            println!("Unable to open database: {}", error);
            return;
        }
    };
}
