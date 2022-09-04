use std::{fs, path::Path};

pub struct Database {
    pub path: String,
}

impl Database {
    pub fn open(path: String) -> Result<Self, String> {
        match fs::create_dir_all(&path) {
            Ok(_) => Ok(Database { path }),
            Err(error) => Err(error.to_string()),
        }
    }

    pub fn set(&self, key: &str, value: &str) -> Result<(), String> {
        println!("{:?}", Path::new(self.path.as_str()).join(key));
        Ok(())
    }
}
