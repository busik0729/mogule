import { Injectable, EventEmitter } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { tap } from 'rxjs/operators';
import { Observable, BehaviorSubject } from 'rxjs';

import { LocalStorageService } from 'angular-web-storage';
import { User } from './models/user';
import { AuthResponse } from './models/auth-response';

import { DeviceInfo } from '../helpers/device_info';
import { Tokens } from "../helpers/tokens";
import { environment } from "../../environments/environment";

@Injectable({
    providedIn: 'root'
})
export class AuthService {

    AUTH_SERVER_ADDRESS = environment.host;
    authSubject = false;

    onUserDataChanged: BehaviorSubject<any>;

    public loginEvent: EventEmitter<boolean>;

    ACCESS_TOKEN = 'ACCESS_TOKEN';
    REFRESH_TOKEN = 'REFRESH_TOKEN';
    EXPIRED_IN = 'EXPIRED_IN';
    USER = 'User';

    redirectUrl: string;

    constructor(
        private httpClient: HttpClient,
        private storage: LocalStorageService,
        private d: DeviceInfo,
        private tkns: Tokens
    ) {
        this.loginEvent = new EventEmitter();
        this.isLoggedIn();
        this.onUserDataChanged = new BehaviorSubject({});
    }

    register(user: User): Observable<AuthResponse> {

        let h = new HttpHeaders();
        h = h.set('X-Device', this.d.getInfoArr());
        const requestOptions = {headers: h};

        return this.httpClient.post<AuthResponse>(`${this.AUTH_SERVER_ADDRESS}user/signup`, user, requestOptions).pipe(
            tap(async (res: AuthResponse) => {
                if (res.Code === 200 && res.Data.Device && res.Data.User) {
                    await this.storage.set(this.ACCESS_TOKEN, res.Data.Device.AccessToken);
                    await this.storage.set(this.REFRESH_TOKEN, res.Data.Device.RefreshToken);
                    await this.storage.set(this.EXPIRED_IN, res.Data.Device.ExpiredIn);
                    await this.storage.set(this.USER, res.Data.User);
                    this.authSubject = true;
                    this.loginEvent.emit(true);
                } else {
                    this.authSubject = false;
                    this.loginEvent.emit(false);
                }
            })
        );
    }

    login(user: User): Observable<AuthResponse> {
        let h = new HttpHeaders();
        h = h.set('X-Device', this.d.getInfoArr());
        const requestOptions = {headers: h};

        return this.httpClient.post(`${this.AUTH_SERVER_ADDRESS}user/login`, user, requestOptions).pipe(
            tap(async (res: AuthResponse) => {

                if (res.Code === 200 && res.Data.Device && res.Data.User) {
                    await this.storage.set(this.ACCESS_TOKEN, res.Data.Device.AccessToken);
                    await this.storage.set(this.REFRESH_TOKEN, res.Data.Device.RefreshToken);
                    await this.storage.set(this.EXPIRED_IN, res.Data.Device.ExpiredIn);
                    await this.storage.set(this.USER, res.Data.User);
                    this.authSubject = true;
                    this.loginEvent.emit(true);
                } else {
                    this.authSubject = false;
                    this.loginEvent.emit(false);
                }
            })
        );
    }

    async logout() {
        await this.storage.remove(this.ACCESS_TOKEN);
        await this.storage.remove(this.REFRESH_TOKEN);
        await this.storage.remove(this.EXPIRED_IN);
        await this.storage.remove(this.USER);
        this.authSubject = false;
        this.loginEvent.emit(false);
    }

    proccessIsLoggedIn() {
        const exp = this.storage.get(this.EXPIRED_IN);
        const AT = this.storage.get(this.ACCESS_TOKEN);
        const RT = this.storage.get(this.REFRESH_TOKEN);
        const user = this.storage.get(this.USER);
        if (exp && AT && RT && user) {
            const n = (new Date().getTime() / 1000) << 0;
            if (exp < n) {
                // this.refreshToken().subscribe()
            }

            this.authSubject = true;
        } else {
            this.authSubject = false;
        }
    }

    isLoggedIn() {
        this.proccessIsLoggedIn();

        return this.authSubject;
    }

    refreshToken(): Observable<string> {
        if (!this.isNeedRefreshToken()) {
            return new Observable<string>();
        }

        let h = new HttpHeaders();
        h = h.set('X-Device', this.d.getInfoArr());
        h = h.set('X-RT', this.tkns.getRT());
        const requestOptions = {headers: h};

        return this.httpClient.get<string>(`${this.AUTH_SERVER_ADDRESS}refresh-token`, requestOptions)
    }

    isNeedRefreshToken(): boolean {
        //expires_at - время когда токен должен истечь, записано при логине или после очередного рефреша
        let expiresAtString = this.storage.get(this.EXPIRED_IN);
        if (!expiresAtString) {
            return false;
        }

        const expiresAt = JSON.parse(expiresAtString);
        //считаем, что токен нужно рефрешить не когда он уже истек, а за минуту до его невалидности
        let isExpireInMinute = new Date().getTime() > (expiresAt - 60);
        return isExpireInMinute;
    }

    refreshStorage(res) {
        this.storage.set(this.ACCESS_TOKEN, res.Data.AccessToken);
        this.storage.set(this.REFRESH_TOKEN, res.Data.RefreshToken);
        this.storage.set(this.EXPIRED_IN, res.Data.ExpiredIn);
    }

    getUserData() {
        return this.storage.get(this.USER)
    }

    setUserData(user) {
        this.storage.set(this.USER, user);
        this.onUserDataChanged.next(user);
    }
}
