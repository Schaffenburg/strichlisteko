export default class Transaction{
    public id: number;
    public value: number;
    public product: string;
    public userId: number;
    public time: Date;
    public undone: boolean;

    constructor(id:number, value:number, product:string, userId:number, time:Date, undone:boolean = false) {
        this.id=id;
        this.value=value;
        this.product=product;
        this.userId=userId;
        this.time=time;
        this.undone=undone;
    }

    /**
     * Gets all transactions for a user
     * @param userId
     */
    public static async getForUser(userId:number): Promise<Transaction[]>{
        let transactions:Transaction[] = []

        return transactions;
    }

    /**
     * Undoes a transaction
     */
    public async undo(){

    }

}
