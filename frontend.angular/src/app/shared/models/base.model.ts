export abstract class BaseModel {
  public id: number;
  public name: string;

  public toString() {
    if (this.id === -1) {
      return this.name + ' (*)';
    }
    return this.name;
  }
}
