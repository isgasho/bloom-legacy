pub const SNIFF_LENGTH: usize = 512;


/// Returns whether a buffer is JPEG image data.
pub fn is_jpeg(buf: &[u8]) -> bool {
    return buf.len() > 2
        && buf[0] == 0xFF
        && buf[1] == 0xD8
        && buf[2] == 0xFF;
}

/// Returns whether a buffer is PNG image data.
pub fn is_png(buf: &[u8]) -> bool {
    return buf.len() > 7
        && buf[0] == 0x89
        && buf[1] == 0x50
        && buf[2] == 0x4E
        && buf[3] == 0x47
        && buf[4] == 0x0D
        && buf[5] == 0x0A
        && buf[6] == 0x1A
        && buf[7] == 0x0A;
}


/// detect_content_type implements the algorithm described at https://mimesniff.spec.whatwg.org/
/// to determine the Content-Type of the given data. It considers at most the first 512 bytes of data.
/// detect_content_type always returns a valid MIME type: if it cannot determine a more specific one,
/// it returns "application/octet-stream".
/// It currently supports only JPEG and PNG.
pub fn detect_content_type(data: &[u8]) -> mime::Mime {
    if is_jpeg(data) {
        return mime::IMAGE_JPEG;
    } else if is_png(data) {
        return mime::IMAGE_PNG;
    }
    return mime::APPLICATION_OCTET_STREAM;
}


#[cfg(test)]
mod tests {
    #[test]
    fn is_jpeg() {
        let v1: Vec<u8> = vec!(1, 2, 3);
        assert_eq!(super::is_jpeg(&v1), false);
        let v2: Vec<u8> = vec!(0xFF, 0xD8, 0xFF);
        assert_eq!(super::is_jpeg(&v2), true);
        let v3: Vec<u8> = vec!(0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A);
        assert_eq!(super::is_jpeg(&v3), false);
    }

    #[test]
    fn is_png() {
        let v1: Vec<u8> = vec!(1, 2, 3);
        assert_eq!(super::is_png(&v1), false);
        let v2: Vec<u8> = vec!(0xFF, 0xD8, 0xFF);
        assert_eq!(super::is_png(&v2), false);
        let v3: Vec<u8> = vec!(0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A);
        assert_eq!(super::is_png(&v3), true);
    }


    #[test]
    fn detect_content_type() {
        let v1: Vec<u8> = vec!(1, 2, 3);
        assert_eq!(super::detect_content_type(&v1), mime::APPLICATION_OCTET_STREAM);
        let v2: Vec<u8> = vec!(0xFF, 0xD8, 0xFF);
        assert_eq!(super::detect_content_type(&v2), mime::IMAGE_JPEG);
        let v3: Vec<u8> = vec!(0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A);
        assert_eq!(super::detect_content_type(&v3), mime::IMAGE_PNG);
    }
}
