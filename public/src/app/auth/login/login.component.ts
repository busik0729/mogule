import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

import { FuseConfigService } from '@fuse/services/config.service';
import { fuseAnimations } from '@fuse/animations';

import { AuthService } from "../auth.service";
import { Router } from '@angular/router';

import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
    selector: 'login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.scss'],
    encapsulation: ViewEncapsulation.None,
    animations: fuseAnimations
})
export class LoginComponent implements OnInit {
    loginForm: FormGroup;

    /**
     * Constructor
     *
     * @param {FuseConfigService} _fuseConfigService
     * @param {FormBuilder} _formBuilder
     * @param {AuthService} _auth
     * @param {Router} _router
     * @param {MatSnackBar} _snackBar
     */
    constructor(
        private _fuseConfigService: FuseConfigService,
        private _formBuilder: FormBuilder,
        private _auth: AuthService,
        private _router: Router,
        private _snackBar: MatSnackBar
    ) {
        // Configure the layout
        this._fuseConfigService.config = {
            layout: {
                navbar: {
                    hidden: true
                },
                toolbar: {
                    hidden: true
                },
                footer: {
                    hidden: true
                },
                sidepanel: {
                    hidden: true
                }
            }
        };
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Lifecycle hooks
    // -----------------------------------------------------------------------------------------------------

    /**
     * On init
     */
    ngOnInit(): void {
        this.loginForm = this._formBuilder.group({
            username: ['', [Validators.required]],
            password: ['', Validators.required]
        });

        if (this._auth.isLoggedIn()) {
            this._router.navigateByUrl('/sample');
        }

        this._auth.loginEvent.subscribe(b => {
            if (b) {
                this._router.navigateByUrl('/sample');
            }
        })
    }

    login(form) {
        this._auth.login(form.value).subscribe((res) => {
            if (res.Code === 200 && res.Data.User && res.Data.Device) {
                this._auth.authSubject = true;

                this._fuseConfigService.config = {
                    layout: {
                        navbar: {
                            hidden: false
                        },
                        toolbar: {
                            hidden: false
                        },
                        footer: {
                            hidden: false
                        },
                        sidepanel: {
                            hidden: false
                        }
                    }
                };

                return false;
            }
        }, (err) => {
            this._snackBar.open(err.error.Error.Message, "OK", {
                duration: 0,
                verticalPosition: "top"
            });
        });
    }
}
