pub mod config {
    pub const TLDS: [&str; 7] = [
        ".com", ".net", ".org", ".lol", ".xyz", ".gg", ".bio"
    ];

    pub struct Database {
        pub mongodb: String,
        pub redis: String,
    }
} 