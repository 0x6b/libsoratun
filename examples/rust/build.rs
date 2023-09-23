fn main() {
    #[cfg(target_os = "macos")]
    {
        let path = "../../lib/archive";
        let lib = "soratun";

        println!("cargo:rustc-link-search=native={}", path);
        println!("cargo:rustc-link-lib=static={}", lib);
    }

    #[cfg(not(target_os = "macos"))]
    {
        todo!()
    }
}
