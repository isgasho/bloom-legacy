package users

/*
fn username_to_salt(username: &str) -> argon2id::Salt {
    let padded_username = pad_or_cut_str(username, argon2id::SALTBYTES);

    let mut salt = argon2id::Salt([0; argon2id::SALTBYTES]);
    (salt.0).clone_from_slice(padded_username.as_bytes());
    salt
}

fn pad_or_cut_str(s: &str, size: usize) -> String {
    let len = s.len();

    match len.cmp(&size) {
        Ordering::Equal => s.to_owned(),
        Ordering::Less => {
            let mut out = String::new();
            out.push_str(s);
            for _ in 0..size - len {
                out.push(0x0 as char);
            }
            out
        }
        Ordering::Greater => s[..size].to_owned(),
    }
}

fn derive_auth_key(username: &str, password: &str) -> Result<Vec<u8>, BloomError> {
    let client_salt = username_to_salt(username);

    let pw_key = argon2id::derive_from_password(
        KDF_PW_KEYLEN,
        password.as_bytes(),
        &client_salt,
        argon2id::OPSLIMIT_INTERACTIVE,
        argon2id::MEMLIMIT_INTERACTIVE,
    )?;

    let auth_key = blake2b::derive_from_key(
        KDF_AUTH_KEYLEN,
        KDF_AUTH_ID,
        KDF_AUTH_CONTEXT,
        &pw_key.as_slice().into(),
    )?;

    return Ok(auth_key);
}

#[cfg(test)]
mod test {
    #[test]
    fn pad_or_cut_str() {
        assert_eq!("", super::pad_or_cut_str("hello world", 0));
        assert_eq!("hello", super::pad_or_cut_str("hello world", 5));
        assert_eq!("hello world", super::pad_or_cut_str("hello world", 11));
        assert_eq!(
            "hello world\x00\x00",
            super::pad_or_cut_str("hello world", 13)
        );
    }
}


// example

fn main() {
    let input = bloom_messages::auth::RegistrationStart {
        display_name: "some display name".to_string(),
        email: "some.email@protonmail.com".to_string(),
    };

    let res = bloom_auth::registration_start(input);
    println!("res = {:?}", res);
}

*/
