import { Injectable } from '@angular/core';
import { LocalStorageService } from "angular-web-storage";
import * as fingerprint2 from 'fingerprintjs2';

@Injectable({
    providedIn: 'root'
})
export class DeviceInfo {
    UUID: string;
    BrowserInfo: any;

    VersionApp = "0.1";

    constructor(private _storage: LocalStorageService) {
        this.init();
    }

    getInfo() {

        return {
            platform: this.BrowserInfo.userAgent,
            uuid: this.getBrowserUUID(),
            model: this.BrowserInfo.platform,
            serial:"serial",
            versionOS: "versionOS",
            versionApp: this.VersionApp
        };
    }
    getInfoArr() {
        return [
            this.BrowserInfo.userAgent,
            this.getBrowserUUID(),
            this.BrowserInfo.platform,
            "serial",
            "versionOS",
            this.VersionApp
        ].join("!");
    }

    getBrowserUUID() {
        return this.UUID;
    }

    getBrowserInfo() {
        return this.BrowserInfo;
    }

    init() {
        let BI = this._storage.get("BI");
        let UUID = this._storage.get("BIUUID");

        if (BI == undefined || UUID == undefined) {
            fingerprint2.get((c) => {
                var b = {};
                var values = c.map(function (component) {
                    b[component.key] = component.value;
                    return component.value
                });
                this.BrowserInfo = b;
                this.UUID = fingerprint2.x64hash128(values.join(''), 31);

                this._storage.set("BI", this.BrowserInfo);
                this._storage.set("BIUUID", this.UUID);
            })
        } else {
            this.BrowserInfo = BI;
            this.UUID = UUID;
        }
    }
}
