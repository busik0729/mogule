import { Injectable, OnDestroy, Inject } from '@angular/core';
import { Observable, SubscriptionLike, Subject, Observer, interval } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { WebSocketSubject, WebSocketSubjectConfig } from 'rxjs/webSocket';

import { share, distinctUntilChanged, takeWhile } from 'rxjs/operators';
import { IWebsocketService, IWsMessage, WebSocketConfig } from './websocket.interfaces';
import { config } from './websocket.config';
import { environment } from "../../environments/environment";

import { Tokens } from "../helpers/tokens";
import { AuthService} from "../auth/auth.service";
import { MatSnackBar } from "@angular/material";

@Injectable()
export class WebsocketService implements IWebsocketService, OnDestroy {
    private config: WebSocketSubjectConfig<IWsMessage<any>>;

    private websocketSub: SubscriptionLike;
    private statusSub: SubscriptionLike;

    private reconnection$: Observable<number>;
    private websocket$: WebSocketSubject<IWsMessage<any>>;
    private connection$: Observer<boolean>;
    private wsMessages$: Subject<IWsMessage<any>>;

    private reconnectInterval: number;
    private reconnectAttempts: number;
    private isConnected: boolean;
    private eventsCount: number;


    public status: Observable<boolean>;

    constructor(
        @Inject(config) private wsConfig: WebSocketConfig,
        private _tkns: Tokens,
        private _auth: AuthService,
        private _snackBar: MatSnackBar,
    ) {
        this.wsMessages$ = new Subject<IWsMessage<any>>();

        this.reconnectInterval = wsConfig.reconnectInterval || 5000; // pause between connections
        this.reconnectAttempts = wsConfig.reconnectAttempts || 50; // number of connection attempts

        this.loadConfig();

        // connection status
        this.status = new Observable<boolean>((observer) => {
            this.connection$ = observer;
        }).pipe(share(), distinctUntilChanged());

        // run reconnect if not connection
        this.statusSub = this.status
            .subscribe((isConnected) => {
                this.isConnected = isConnected;

                if (!this.reconnection$ && typeof (isConnected) === 'boolean' && !isConnected && this._auth.isLoggedIn()) {
                    console.log("reconn");
                    this.reconnect();
                }
            });

        this.websocketSub = this.wsMessages$.subscribe(
            null, (error: ErrorEvent) => console.error('WebSocket error!', error)
        );

        if (this._auth.isLoggedIn()) {
            this.connect();
        }

        this._auth.loginEvent.subscribe(b => {
            if (b) {
                this.connect();
            } else {
                this.close();
            }
        });

        this.eventsCount = 0;
    }

    ngOnDestroy() {
        this.websocketSub.unsubscribe();
        this.statusSub.unsubscribe();
    }

    private loadConfig() {
        this.config = {
            url: environment.ws + "?X-AT=" + this._tkns.getAT(),
            closeObserver: {
                next: (event: CloseEvent) => {
                    console.log("close observer");
                    this.websocket$ = null;
                    this.connection$.next(false);
                }
            },
            openObserver: {
                next: (event: Event) => {
                    console.log('WebSocket connected!');
                    this.connection$.next(true);
                }
            }
        };
    }


    /*
    * connect to WebSocked
    * */
    private connect(): void {

        this.loadConfig();

        this.websocket$ = new WebSocketSubject(this.config);

        this.websocket$.subscribe(
            (message) => {
                this.wsMessages$.next(message)
            },
            (error: Event) => {
                if ((this.isConnected == undefined || !this.isConnected) && this._auth.isLoggedIn()) {
                    // run reconnect if errors

                    this._auth.refreshToken().subscribe(authHeader => {
                        this._auth.refreshStorage(authHeader);
                        // this._auth.loginEvent.emit(true);
                        this.reconnect();
                    }, (e) => {
                        this.close();
                        this._auth.logout();
                    });
                    // this.reconnect();
                }
            });
    }

    private close(): void {
        this.websocket$.unsubscribe();
        this.connection$.next(false);

    }


    /*
    * reconnect if not connecting or errors
    * */
    private reconnect(): void {
        this.reconnection$ = interval(this.reconnectInterval)
            .pipe(takeWhile((v, index) => index < this.reconnectAttempts && !this.websocket$));

        this.reconnection$.subscribe(
            () => this.connect(),
            null,
            () => {
                // Subject complete if reconnect attemts ending
                this.reconnection$ = null;

                if (!this.websocket$) {
                    this.wsMessages$.complete();
                    this.connection$.complete();
                }
            });
    }

    public incrEventsCount() {
        this.eventsCount++;
    }


    /*
    * on message event
    * */
    public on<T>(event: string): Observable<T> {
        if (event) {
            // this._snackBar.open("У вас " + this.eventsCount + " уведомлений", "OK", {
            //     duration: 5000,
            //     verticalPosition: "top"
            // });

            return this.wsMessages$.pipe(
                filter((message: IWsMessage<T>) => message.event === event),
                map((message: IWsMessage<T>) => message.data)
            );
        }
    }


    /*
    * on message to server
    * */
    public send(event: string, data: any = {}): void {
        if (event && this.isConnected && this.websocket$) {
            console.log(JSON.stringify({ event, data }));
            this.websocket$.next({ event, data });
        } else {
            console.error('Send error!');
        }
    }
}
