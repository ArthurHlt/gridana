import {Injectable} from '@angular/core';
import {PersistenceService, StorageType} from 'angular-persistence';

@Injectable()
export class ConfigService {

  constructor(private persistenceService: PersistenceService) {

  }

  defineShowSilenced(obj: any, propName: string) {
    this.persistenceService.defineProperty(obj, propName, 'show_silenced', {type: StorageType.LOCAL});
  }

  defineReceiveResolved(obj: any, propName: string) {
    this.persistenceService.defineProperty(obj, propName, 'receive_resolved', {type: StorageType.LOCAL});
  }
}
