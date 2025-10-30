module.exports = {
    printWidth: 120,
    useTabs: false,
    tabWidth: 4,
    singleQuote: true,
    trailingComma: 'none',
    bracketSpacing: true,
    endOfLine: 'auto',
    quoteProps: 'as-needed',
    bracketSameLine: false,
    htmlWhitespaceSensitivity: 'ignore',
    overrides: [
        {
            files: '*.html',
            options: {
                parser: 'angular'
            },
        },
        {
            files: '*.scss',
            options: {
                parser: 'scss',
                singleQuote: true,
            },
        },
        {
            files: ['*.yaml', '*.yml'],
            options: {
                parser: 'yaml',
                singleQuote: false,
            }
        }
    ]
};
