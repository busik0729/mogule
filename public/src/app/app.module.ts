import { Injector, NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterModule, Routes } from '@angular/router';
import { MatMomentDateModule } from '@angular/material-moment-adapter';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { TranslateModule } from '@ngx-translate/core';
import 'hammerjs';

import { FuseModule } from '@fuse/fuse.module';
import { FuseSharedModule } from '@fuse/shared.module';
import { FuseProgressBarModule, FuseSidebarModule, FuseThemeOptionsModule } from '@fuse/components';

import { fuseConfig } from 'app/fuse-config';

import { AppComponent } from 'app/app.component';
import { LayoutModule } from 'app/layout/layout.module';
import { SampleModule } from 'app/main/sample/sample.module';
import { ScrumboardModule } from 'app/main/scrumboard/scrumboard.module';
import { ContactsModule } from "./main/contacts/contacts.module";
import { ProfileModule } from "./main/profile/profile.module";
import { ContactsModule as UsersModule } from "./main/users/contacts.module";

import { WebsocketModule } from './websocket';

import { AngularWebStorageModule } from 'angular-web-storage';
import { AuthModule } from './auth/auth.module';

import { HTTP_INTERCEPTORS } from "@angular/common/http";
import { AuthenticationInterceptor } from "./interceptors/auth.interceptor";
import { AuthService } from "./auth/auth.service";
import { DeviceInfo } from "./helpers/device_info";

import { InMemoryWebApiModule } from 'angular-in-memory-web-api';
import { FakeDbService } from "./fake-db/fake-db.service";
import { environment } from "../environments/environment";

const appRoutes: Routes = [
    {
        path: '**',
        redirectTo: 'sample'
    },
    {
        path: 'auth',
        loadChildren: './auth/auth.module#AuthModule'
    },
    {
        path: 'scrum',
        loadChildren: './main/scrumboard/scrumboard.module#ScrumboardModule'
    },
    {
        path: 'contacts',
        loadChildren: './main/contacts/contacts.module#ContactsModule'
    },
    {
        path: 'users',
        loadChildren: './main/users/contacts.module#ContactsModule'
    },
    {
        path: 'profile',
        loadChildren: './main/profile/profile.module#ProfileModule'
    },
];

@NgModule({
    declarations: [
        AppComponent
    ],
    imports: [
        BrowserModule,
        BrowserAnimationsModule,
        HttpClientModule,
        RouterModule.forRoot(appRoutes),

        TranslateModule.forRoot(),

        InMemoryWebApiModule.forRoot(FakeDbService, {
            delay: 0,
            passThruUnknownUrl: true
        }),

        // Material moment date module
        MatMomentDateModule,

        // Material
        MatButtonModule,
        MatIconModule,

        // Fuse modules
        FuseModule.forRoot(fuseConfig),
        FuseProgressBarModule,
        FuseSharedModule,
        FuseSidebarModule,
        FuseThemeOptionsModule,

        // App modules
        LayoutModule,
        SampleModule,
        ScrumboardModule,
        ContactsModule,
        UsersModule,
        ProfileModule,
        // AnalyticsDashboardModule,
        // CalendarModule,

        AngularWebStorageModule,
        AuthModule,

        WebsocketModule.config({
            url: environment.ws
        })
    ],
    providers: [
        {
            provide: HTTP_INTERCEPTORS,
            useClass: AuthenticationInterceptor,
            multi: true,
        }
    ],
    bootstrap: [
        AppComponent
    ]
})
export class AppModule {
    constructor(inj: Injector,
                auth: AuthService,
                device: DeviceInfo,
                http: HttpClient) {
        device.init();
        //Получаем интерцепторы которые реализуют интерфейс AuthInterceptor
        let interceptors = inj.get<AuthenticationInterceptor[]>(HTTP_INTERCEPTORS)
            .filter(i => {
                return i.init;
            });
        //передаем http сервис и сервис авторизации.
        interceptors.forEach(i => i.init(http, auth));
    }
}
