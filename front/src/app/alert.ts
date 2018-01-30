export class Alert {
  id: string;
  name: string;
  startsAt: string;
  endsAt: string;
  generatorURL: string;
  notifierURL: string;
  status: string;
  labels: Map<string, string>;
  annotations: Map<string, string>;
  color: string;
  message: string;
  probe: string;
  identifier: string;
  weight: number;
}