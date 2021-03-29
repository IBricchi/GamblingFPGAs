// synopsys translate_off
'timescale 1 ps / 1 ps
// synopsys translate_on

module hex_to_7segs (out,in,whichSeg,clear, clock);
	output [41:0] out;
	input [3:0] in;
	input [5:0] whichSeg;
	input clear;
	reg [31:0] out;
	reg [6:0] calc;
	input	clock;


	/*always @ (*)
		case (in)
			4'h0: calc = 7'b1000000;
			4'h1: calc = 7'b1111001;	//   -- 0 ---
			4'h2: calc = 7'b0100100; 	// |		 |
			4'h3: calc = 7'b0110000; 	// 5         1
			4'h4: calc = 7'b0011001; 	// |	 	 |
			4'h5: calc = 7'b0010010; 	//   -- 6 ---
			4'h6: calc = 7'b0000010; 	// |	 	 |
			4'h7: calc = 7'b1111000; 	// 4         2
			4'h8: calc = 7'b0000000; 	// |	 	 |
			4'h9: calc = 7'b0011000;
			default: calc=7'b0000000;//   -- 3 ---
		endcase
		*/
	always @ (whichSeg or clear) begin
		if(clear) begin
		//out[41:0] <= 42'b100000010000001000000100000010000001000000;
		end
		else begin
			case (whichSeg)
				6'h1: out[3:0]   <= in;
				6'h2: out[7:4]  <= in;
				6'h4: out[11:8] <= in;
				6'h8: out[15:12] <= in;
				6'h10: out[19:16] <= in;
				6'h20: out[23:20] <= in;
			endcase
		end
		end
	endmodule
