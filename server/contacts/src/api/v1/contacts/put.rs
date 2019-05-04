use kernel::{
    api,
    log::macros::*,
    api::middlewares::{
        GetRequestLogger,
        GetRequestId,
        GetRequestAuth,
    },
    KernelError,
};
use futures::{
    future::{
        Either,
        ok,
        Future,
    },
};
use actix_web::{
    web, Error, HttpRequest, HttpResponse, ResponseError,
};
use crate::{
    api::v1::models,
    controllers,
};


pub fn put(contact_id: web::Path<(uuid::Uuid)>, contact_data: web::Json<models::UpdateContactBody>, state: web::Data<api::State>, req: HttpRequest)
-> impl Future<Item = HttpResponse, Error = Error> {
    let logger = req.logger();
    let auth = req.request_auth();
    let request_id = req.request_id().0;

    if auth.session.is_none() || auth.account.is_none() {
        return Either::A(ok(KernelError::Unauthorized("Authentication required".to_string()).error_response()));
    }

    return Either::B(
        state.db.send(controllers::UpdateContact{
            addresses: contact_data.addresses.clone(),
            birthday: contact_data.birthday.clone(),
            company: contact_data.company.clone(),
            emails: contact_data.emails.clone(),
            first_name: contact_data.first_name.clone(),
            last_name: contact_data.last_name.clone(),
            notes: contact_data.notes.clone(),
            occupation: contact_data.occupation.clone(),
            organizations: contact_data.organizations.clone(),
            phones: contact_data.phones.clone(),
            websites: contact_data.websites.clone(),
            contact_id: contact_id.into_inner(),
            actor_id: auth.account.expect("error unwraping non none account").id,
            session_id: auth.session.expect("error unwraping non none session").id,
            request_id,
        })
        .map_err(|_| KernelError::ActixMailbox)
        .from_err()
        .and_then(move |contact| {
            match contact {
                Ok(contact) => {
                    let res = models::ContactResponse{
                        id: contact.id,
                        created_at: contact.created_at,
                        updated_at: contact.updated_at,
                        addresses: contact.addresses,
                        birthday: contact.birthday,
                        company: contact.company,
                        emails: contact.emails,
                        first_name: contact.first_name,
                        last_name: contact.last_name,
                        notes: contact.notes,
                        occupation: contact.occupation,
                        organizations: contact.organizations,
                        phones: contact.phones,
                        websites: contact.websites,
                    };
                    let res = api::Response::data(res);
                    ok(HttpResponse::Ok().json(&res))
                },
                Err(err) => {
                    slog_error!(logger, "{}", err);
                    ok(err.error_response())
                },
            }
        })
    );
}
