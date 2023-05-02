fn main() {
    let path = "../lib/archive";
    let lib = "soratun";

    println!("cargo:rustc-link-search=native={}", path);
    println!("cargo:rustc-link-lib=static={}", lib);
}
