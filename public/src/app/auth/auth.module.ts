import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { AngularWebStorageModule } from 'angular-web-storage';

import { LoginModule } from "./login/login.module";
import { RegisterModule } from "./register/register.module";

@NgModule({
    declarations: [],
    imports: [
        CommonModule,
        HttpClientModule,
        AngularWebStorageModule,
        LoginModule,
        RegisterModule
    ]
})
export class AuthModule {
}
