// Path to org1 user private key directory.
export type Order = {
    bidMatchId: number;
    bidStatus: number;
    id: number;
    onMarketPrice: number;
    orderCost: number;
    paymentId: number;
    slotId: string;
    slotExecDate: number;
    totalQuantity: number;
    unitCost: number;
    action: number;
    userId: number;
};

// Path to org1 user private key directory.
export type Payment = {
    paymentID: number;
    paymentType: string;
    totalAmount: number;
    paymentDetailID: number;
    debitedFrom: number;
    creditedTo: number;
    totalUnitCost: number;
    platformFee: number;
    tokenAmount: number;
    bidRefundAmount: number;
    platformFeeRefundAmount: number;
    penaltyFromSeller: number;
};
