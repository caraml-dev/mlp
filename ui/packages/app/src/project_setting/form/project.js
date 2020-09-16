const objectAssignDeep = require(`object-assign-deep`);

export class Project {
  constructor() {
    this.name = "";
    this.team = "";
    this.stream = "";
    this.administrators = [];
    this.readers = [];
    this.labels = [];
  }

  static from(json) {
    return objectAssignDeep(new Project(), json);
  }
}
