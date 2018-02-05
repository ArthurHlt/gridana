import {Alert} from "./alert";

export class OrderedAlerts {
  alerts: Map<string, Map<string, Alert[]>>;
  identifiers: string[];
}