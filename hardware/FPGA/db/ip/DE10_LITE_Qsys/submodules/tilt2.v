// synopsys translate_off
'timescale 1 ps / 1 ps
// synopsys translate_on
module tilt (
	aclr,
	clk_en,
	clock,
	dataa,
	datab,
	datac,
	result
);

	input	aclr;
	input	clk_en;
	input	clock;
	input	[31:0]	dataa;
	input	[31:0]	datab;
	input		datac;
	output	[1:0]	result;

	reg [31:0] result; // = dataa[15:0] << datab[15:0];
	reg [31:0] all;
	reg [31:0] a;
	always @ (dataa) begin
	 	all <= ((8'd100 - datab)) + 12'd400;
		a <= dataa * 4'd10;

		if(datac == 1'b1)
			result <= 2'b00;
		else begin
			if(a < 16'h15e)
				result <= 2'b00;
			else begin
				if(a >= 12'd350 && a < all)
					result <= 2'b01;
				else 
					result <= 2'b11;
			end
		end
		
	end

endmodule
