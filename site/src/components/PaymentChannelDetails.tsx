type Props = {
  paymentChannel: string;
};

export default function PaymentChannelDetails({ paymentChannel }: Props) {
  return <p>{paymentChannel}</p>;
}
