// synopsys translate_off
'timescale 1 ps / 1 ps
// synopsys translate_on
module fir (
	aclr,
	clk_en,
	clock,
	dataa,
	datab,
	result
);

	input	aclr;
	input	clk_en;
	input	clock;
	input	[31:0]	dataa;
	input	[31:0]	datab;
	output	[31:0]	result;

	reg 	[31:0] 	result; // = dataa[15:0] << datab[15:0];
	
	reg	 	newClk;
	reg 	[14:0] 	count; 

    reg [31:0] dataInEx;
	reg [31:0] stage1;
	reg [31:0] stage2;
	reg [31:0] stage3;
	reg [31:0] stage4;
	reg [31:0] stage5;
	reg [31:0] stage6;
	reg [31:0] stage7;
	reg [31:0] stage8;
	reg [31:0] stage9;
	reg [31:0] stage10;
	reg [31:0] stage11;
	reg [31:0] stage12;
	reg [31:0] stage13;
	reg [31:0] stage14;
	reg [31:0] stage15;
	reg [31:0] stage16;
	reg [31:0] stage17;
	reg [31:0] stage18;
	reg [31:0] stage19;
	reg [31:0] stage20;
	reg [31:0] stage21;
	reg [31:0] stage22;
	reg [31:0] stage23;
	reg [31:0] stage24;
	reg [31:0] stage25;
	reg [31:0] stage26;
	reg [31:0] stage27;
	reg [31:0] stage28;
	reg [31:0] stage29;
	reg [31:0] stage30;
	reg [31:0] stage31;
	reg [31:0] stage32;
	reg [31:0] stage33;
	reg [31:0] stage34;
	reg [31:0] stage35;
	reg [31:0] stage36;
	reg [31:0] stage37;
	reg [31:0] stage38;
	reg [31:0] stage39;
	reg [31:0] stage40;
	reg [31:0] stage41;
	reg [31:0] stage42;
	reg [31:0] stage43;
	reg [31:0] stage44;
	reg [31:0] stage45;
	reg [31:0] stage46;
	reg [31:0] stage47;
	reg [31:0] stage48;
	reg [31:0] stage49;



	always @(posedge clock) begin
		count <= count + 1;
		if(count == 14'd12499) begin
			newClk <= newClk + 1;
			count <= 0;
		end
	end


	always @ (newClk) begin
		
		stage1 	<= dataa;
		stage2 	<= stage1;
			
		
result <= (2*dataa) + (2*stage1);

		
	end

endmodule
