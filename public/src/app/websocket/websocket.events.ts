export const WS = {
    ON: {
        MESSAGES: 'message',
        COUNTER: 'counter',
        UPDATE_TEXTS: 'update-texts',
        BOARD: {
            CREATE: 'create-board',
            UPDATE: 'update-board',
            DELETE: 'delete-board',
            DEADLINE: 'deadline-board'
        },
        LIST: {
            CREATE: 'create-list',
            UPDATE: 'update-list',
            DELETE: 'delete-list'
        },
        CARD: {
            CREATE: 'create-card',
            UPDATE: 'update-card',
            DELETE: 'delete-card',
            DEADLINE: 'deadline-card'
        },
        COMMENT: {}
    },
    SEND: {
        SEND_TEXT: 'message',
        REMOVE_TEXT: 'remove-text',
        READ_EVENT: 'event-read'
    }
};
