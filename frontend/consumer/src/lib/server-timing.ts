export const ServerTimingHeaderName = 'Server-Timing';

export class ServerTimingEvent {
  name = '';
  description = '';
  startTime: Date = new Date();
  endTime: Date = new Date();
  duration = 0;

  constructor(name = '', description = '') {
    this.name = name;
    this.description = description;
    this.startTime = new Date();
  }

  end(): void {
    this.endTime = new Date();
    this.duration = this.endTime.getTime() - this.startTime.getTime();
  }
}

export class ServerTiming {
  events: ServerTimingEvent[] = [];

  addEvent(name = '', description = ''): ServerTimingEvent {
    const event = new ServerTimingEvent(name, description);
    this.events.push(event);
    return event;
  }

  headerValue(): string {
    let result = '';
    for (const event of this.events) {
      result += `${event.name};desc="${event.description}";dur=${event.duration},`;
    }
    return result.replace(/,(\s*)$/, '$1');
  }
}
