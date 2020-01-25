import { Injectable } from '@angular/core';
import { LocalStorageService } from 'angular-web-storage';

@Injectable({
    providedIn: 'root'
})
export class Tokens {
    AT: string = "ACCESS_TOKEN";
    RT: string = "REFRESH_TOKEN";
    constructor(private storage: LocalStorageService) {}

    getAT() {
        const at = this.storage.get(this.AT); // переделать LocalStorageService для проверки на существование токенов
        return at;
    }

    getRT() {
        const rt = this.storage.get(this.RT); // переделать LocalStorageService для проверки на существование токенов
        return rt;
    }
}
