import { Observable, Subscriber } from 'rxjs';
import { Injectable } from '@angular/core';
import { HttpEvent, HttpInterceptor, HttpHandler, HttpRequest, HttpClient } from '@angular/common/http';
import { AuthService } from "../auth/auth.service";
import { finalize } from "rxjs/operators";

type CallerRequest = {
    subscriber: Subscriber<any>;
    failedRequest: HttpRequest<any>;
};

@Injectable()
export class AuthenticationInterceptor implements HttpInterceptor {

    private auth: AuthService;
    private http: HttpClient;
    private refreshInProgress: boolean;
    private requests: CallerRequest[] = [];

    init(http: HttpClient, auth: AuthService) {
        this.auth = auth;
        this.http = http;
    }

    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {

        let observable = new Observable<HttpEvent<any>>((subscriber) => {
            //как только вызывающий код сделает подписку мы попадаем сюда и подписываемся на наш HttpRequest
            //тобишь выполняем оригинальный запрос
            let originalRequestSubscription = next.handle(req)
                .subscribe((response) => {
                        //оповещаем в инициатор (success) ответ от сервера
                        subscriber.next(response);
                    },
                    (err) => {
                        if (err.status === 401) {
                            //если споймали 401ую - обрабатываем далее по нашему алгоритму
                            this.handleUnauthorizedError(subscriber, req);
                        } else {
                            //оповещаем об ошибке
                            subscriber.error(err);
                        }
                    },
                    () => {
                        //комплит запроса, отрабатывает finally() инициатора
                        subscriber.complete();
                    });

            return () => {
                // на случай если в вызывающем коде мы сделали отписку от запроса
                // если не сделать отписку и здесь, в dev tools браузера не увидим отмены запросов, т.к инициатор (например Controller) делает отписку от нашего враппера, а не от исходного запроса
                originalRequestSubscription.unsubscribe();
            };
        });

//вернем вызывающему коду Observable, пусть сам решает когда делать подписку.
        return observable;
    }

    private handleUnauthorizedError(subscriber: Subscriber < any >, request: HttpRequest<any>) {

        //запоминаем "401ый" запрос
        this.requests.push({ subscriber, failedRequest: request });
        if(!this.refreshInProgress) {
            //делаем запрос на восстанавливение токена, и установим флаг, дабы следующие "401ые"
            //просто запоминались но не инициировали refresh
            this.refreshInProgress = true;
            this.auth.refreshToken()
                .pipe(finalize(() => {
                    this.refreshInProgress = false;
                }))
                .subscribe((authHeader) => {
                    //если токен рефрешнут успешно, повторим запросы которые накопились пока мы ждали ответ от рефреша
                        this.repeatFailedRequests(authHeader);
                        this.auth.refreshStorage(authHeader);
                        this.auth.loginEvent.emit(true)
                },
                () => {
                    //если по каким - то причинам запрос на рефреш не отработал, то делаем логаут
                    this.auth.logout();
                });
        }
    }

    private repeatFailedRequests(authHeader) {

        this.requests.forEach((c) => {
            //клонируем наш "старый" запрос, с добавлением новенького токена
            const requestWithNewToken = c.failedRequest.clone({
                headers: c.failedRequest.headers.set('X-AT', authHeader.Data.AccessToken)
            });
            //и повторяем (помним с.subscriber - subscriber вызывающего кода)
            this.repeatRequest(requestWithNewToken, c.subscriber);
        });
        this.requests = [];
    }

    private repeatRequest(requestWithNewToken: HttpRequest < any >, subscriber: Subscriber<any>) {

        //и собственно сам процесс переотправки
        this.http.request(requestWithNewToken).subscribe((res) => {
                subscriber.next(res);
            },
            (err) => {
                if (err.status === 401) {
                    // if just refreshed, but for unknown reasons we got 401 again - logout user
                    this.auth.logout();
                }
                subscriber.error(err);
            },
            () => {
                subscriber.complete();
            });
    }
}
