import { FuseNavigation } from '@fuse/types';

export const navigation: FuseNavigation[] = [
    {
        id       : 'applications',
        title    : 'Applications',
        translate: 'NAV.APPLICATIONS',
        type     : 'group',
        children : [
            {
                id       : 'sample',
                title    : 'Главная',
                type     : 'item',
                icon     : 'email',
                url      : '/sample',
                // badge    : {
                //     title    : '25',
                //     bg       : '#F44336',
                //     fg       : '#FFFFFF'
                // }
            },
            {
                id       : 'scrumboard',
                title    : 'Проекты',
                type     : 'item',
                icon     : 'assessment',
                url      : '/boards'
            },
            {
                id       : 'contacts',
                title    : 'Контакты',
                type     : 'item',
                icon     : 'account_box',
                url      : '/contacts'
            },
            {
                id       : 'users',
                title    : 'Пользователи',
                type     : 'item',
                icon     : 'account_box',
                url      : '/users',
                for      : ['Администратор', 'Руководитель']
            },
            // {
            //     id       : 'analytics',
            //     title    : 'Analytics',
            //     type     : 'item',
            //     icon     : 'account_box',
            //     url      : '/analytics'
            // },
            // {
            //     id       : 'calendar',
            //     title    : 'Calendar',
            //     type     : 'item',
            //     icon     : 'account_box',
            //     url      : '/calendar'
            // }
        ]
    }
];
