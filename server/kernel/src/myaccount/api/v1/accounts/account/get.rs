use crate::{
    api,
    log::macros::*,
    api::middlewares::{
        GetRequestLogger,
        GetRequestId,
        GetRequestAuth,
    },
    error::KernelError,
    myaccount::{
        controllers,
        api::v1::models,
    },
};
use futures::{
    future::Future,
    future::ok,
    future::Either,
};
use actix_web::{
    web, Error, HttpRequest, HttpResponse, ResponseError,
};


pub fn get(account_id: web::Path<(uuid::Uuid)>, state: web::Data<api::State>, req: HttpRequest)
-> impl Future<Item = HttpResponse, Error = Error> {
    let logger = req.logger();
    let auth = req.request_auth();

    if auth.session.is_none() || auth.account.is_none() {
        return Either::A(ok(KernelError::Unauthorized("Authentication required".to_string()).error_response()));
    }

    return Either::B(
        state.db
        .send(controllers::FindAccountAdmin{
            actor: auth.account.expect("unwraping non none account"),
            account_id: account_id.into_inner(),
        })
        .map_err(|_| KernelError::ActixMailbox)
        .from_err()
        .and_then(move |res: Result<_, KernelError>| {
            match res {
                Ok(account) => {
                    let account_res: models::AccountResponse = From::from(account);
                    let res = api::Response::data(account_res);
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
