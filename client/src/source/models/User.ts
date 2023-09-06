import Product from "@/source/models/Product";
import Api from "@/source/api";

export default class User {
    public id:number
    public username:string
    public imageId:string = ""
    private balance:number = 0
    public active:boolean = true

    constructor(id:number, username:string, imageId:string = "", balance :number = 0, active :boolean = true) {
        this.id=id;
        this.username=username;
        this.imageId=imageId;
        this.balance=balance;
        this.active=active;
    }

    public getBalance(){
        return this.balance;
    }

    public async buyProductId(productId:number){

    }
    public async buyProduct(product:Product){

    }

    public async deposit(amount:number){

    }

    public async withdraw(amount: number){

    }

    public async transactions(){

    }

    public async undoTransaction(transactionId: number){

    }

    public static async all(){
        const api = new Api();
        const resp = await api.getAllUsers();
        if(resp.data){
            const users: User[] = [];
            for(const user of resp.data){
                users.push(new User(user.id,user.username,user.image,user.balance,user.active));
            }
            return users;
        }
        return []
    }



}
