export default {
    verbose: true,
    testEnvironment: 'node',
    transform: {
        '^.+\\.js$': 'babel-jest'
    },
    moduleFileExtensions: ['js', 'json'],
    testMatch: [
        '**/test/**/*.test.js'
    ],
    collectCoverage: true,
    coverageDirectory: 'coverage',
    coverageReporters: ['text', 'lcov'],
    coverageThreshold: {
        global: {
            branches: 80,
            functions: 80,
            lines: 80,
            statements: 80
        }
    },
    testPathIgnorePatterns: [
        '/node_modules/',
        '/dist/'
    ],
    moduleNameMapper: {
        '^@/(.*)$': '<rootDir>/src/$1'
    }
}; 