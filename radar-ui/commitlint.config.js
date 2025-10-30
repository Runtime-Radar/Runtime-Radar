module.exports = {
    extends: ['@koobiq/commitlint-config'],
    rules: {
        'scope-enum': [
            2,
            'always',
            [
                'app',
                'auth'
            ]
        ]
    }
};
