export default {
    displayName: 'domain-detector',
    preset: '../../../jest.preset.js',
    setupFilesAfterEnv: [
        '<rootDir>/src/test-setup.ts'
    ],
    globals: {
        'ts-jest': {
            stringifyContentPathRegex: '\\.(html|svg)$',
            tsconfig: '<rootDir>/tsconfig.spec.json'
        }
    },
    coverageDirectory: '../../../coverage/libs/domain-detector',
    snapshotSerializers: [
        'jest-preset-angular/build/serializers/no-ng-attributes',
        'jest-preset-angular/build/serializers/ng-snapshot',
        'jest-preset-angular/build/serializers/html-comment'
    ],
    transformIgnorePatterns: [
        'node_modules/(?!.*\\.js$)'
    ],
    transform: {
        '^.+\\.(ts|mjs|js|html)$': 'jest-preset-angular'
    }
};
