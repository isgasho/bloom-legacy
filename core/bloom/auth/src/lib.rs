use crypto42::kdf::{argon2id, blake2b};

pub const API_URL: &str = "http://localhost:8000/";
pub const KDF_CONTEXT: &str = "__auth__";

/// `registration_start` sends
///
/// See https://theguide.bloom.sh/projects/bloom/security/authentication.html#registration for the spec
pub fn registration_start(input: messages::auth::GuiRegistrationStart) -> messages::Message {
    let salt = argon2id::gen_salt();

    let pw_key = argon2id::derive_from_password(
        32,
        input.password.as_bytes(),
        &salt,
        argon2id::OPSLIMIT_INTERACTIVE,
        argon2id::MEMLIMIT_INTERACTIVE,
    )
    .expect("error deriving pw_key from password");

    let auth_key = blake2b::derive_from_key(64, 1, KDF_CONTEXT, &pw_key.as_slice().into())
        .expect("error deriving auth_key from pw_key");

    let message: messages::Message = messages::auth::StartRegistration {
        display_name: input.display_name,
        email: input.email,
        auth_key: base64::encode(&auth_key),
    }
    .into();

    let client = reqwest::Client::new();
    let mut api_res = client
        .post(API_URL)
        .json(&message)
        .send()
        .expect("error posting to API");

    let ret: messages::Message = api_res
        .json()
        .expect("error converting api response back to JSON");

    return ret;
}

pub fn registration_verify(message: messages::auth::RegistrationVerify) -> messages::Message {
    let message: messages::Message = message.into();

    let client = reqwest::Client::new();
    let mut api_res = client
        .post(API_URL)
        .json(&message)
        .send()
        .expect("error posting to API");

    let ret: messages::Message = api_res
        .json()
        .expect("error converting api response back to JSON");

    return ret;
}

pub fn registration_complete(message: messages::auth::RegistrationComplete) -> messages::Message {
    let message: messages::Message = message.into();

    let client = reqwest::Client::new();
    let mut api_res = client
        .post(API_URL)
        .json(&message)
        .send()
        .expect("error posting to API");

    let ret: messages::Message = api_res
        .json()
        .expect("error converting api response back to JSON");

    return ret;
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
