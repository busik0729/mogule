export interface AuthResponse {
    Code: number;
    Data: {
        User: {
            Id: string;
            RoleString: string;
            Surname: string;
            Name: string;
            Username: string;
        },
        Device: {
            AccessToken: string;
            RefreshToken: string;
            ExpiredIn: number;
        }
    };
}
