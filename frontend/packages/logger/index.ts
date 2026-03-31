import { Logger } from 'tslog';

export type LoggerType = Logger<void>;

export function buildServerSideLogger(name: string, pretty = false): Logger<void> {
  return new Logger<void>({
    type: pretty ? 'pretty' : 'json',
    hideLogPositionForProduction: true,
    prettyLogTemplate: '{{yyyy}}.{{mm}}.{{dd}} {{hh}}:{{MM}}:{{ss}}:{{ms}}\t{{logLevelName}}\t',
    stylePrettyLogs: false,
    prettyLogTimeZone: 'UTC',
    prettyLogStyles: {
      logLevelName: {
        '*': ['bold', 'black', 'bgWhiteBright', 'dim'],
        'SILLY': ['bold', 'white'],
        'TRACE': ['bold', 'whiteBright'],
        'DEBUG': ['bold', 'green'],
        'INFO': ['bold', 'blue'],
        'WARN': ['bold', 'yellow'],
        'ERROR': ['bold', 'red'],
        'FATAL': ['bold', 'redBright'],
      },
      dateIsoStr: 'white',
      filePathWithLine: 'white',
      name: ['white', 'bold'],
      nameWithDelimiterPrefix: ['white', 'bold'],
      nameWithDelimiterSuffix: ['white', 'bold'],
      errorName: ['bold', 'bgRedBright', 'whiteBright'],
      fileName: ['yellow'],
    },
  });
}
