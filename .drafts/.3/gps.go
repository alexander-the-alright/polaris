// =============================================================================
// Auth: Alex Celani
// File: gps.go
// Revn: 07-25-2023  3.0
// Func: asynchronously receive messages from another program, print
//
// TODO: add LED
//       add flags
// =============================================================================
// CHANGE LOG
// -----------------------------------------------------------------------------
//*07-13-2023: copied over, removed call to alter()/forward()
//             removed any call to write
//             put read and print in a for loop
//             wrote getFileName() to create a filename based off
//              current time
//             write output to a logfile
//             moved file writing to AFTER checking for fin
//*07-20-2023: added calls to os.Create(), file.Write(), etc
// 07-21-2023: tested filewriting stuff, no avail
// 07-24-2023: added bufio import and buffered file writing
//*07-25-2023: commented bufio stuffs
//
// =============================================================================

package main

import ( 
    "bufio"     // NewWriter, Write, Writer.Flush
    "net"       // ResolveTCPAddr, conn.Write,Read,Close
                // ListenTCP, listener.Accept
    "os"        // Args, Stderr, Exit, Create, Write, Close,
                // WriteString
    "fmt"       // Fprintf, Println
    "strings"   // ToLower, Compare
    "time"      // Now, Format
)


// quick error checking
func check( err error ) {
    // if user checks an error, it better be nil
    if err != nil {
        // put error in string, print in stderr
        fmt.Fprintf( os.Stderr, "Fatal error: %s", err.Error() )
        os.Exit( 2 )        // quit
    }
}


func getFileName() string{
    // capture current time
    t := time.Now()
    // format as YYYYMMDDHHmmss.gps
    return t.Format( "20060102150405" ) + ".gps"
}


// function to handle incoming connections
func handleClient( conn net.Conn ) {
    defer conn.Close()  // barring any error, still close connection

    var buf [512]byte   // declare large byte array, store messages

// debug printing stuff
//    fmt.Println( "handling" )

    // iterate forever to always read over connection
    for {
        // read n bytes from connection into buffer
        n, err := conn.Read( buf[0:] )
        if err != nil { // erroring on read will simply leave the 
            return      // function so it can start again later
        }

        // print recv'd message
        // string() only works on byte SLICES so [:] is required
        // debug print
//        fmt.Println( "recv: ", string( buf[:n] ) )

        // essentially a case statement for conn.Read, Compare makes
        // it hard to do that
        // check to see if buf[:n] is start
        if strings.Compare( string( buf[:n] ), "start" ) == 0 {
            filename := getFileName()   // get name of file
            file, err = os.Create( filename )   // create new file
            check( err )    // make sure file creation worked
            writer = bufio.NewWriter( file )    // create filewriter
        // check to see if buf[:n] is fin
        } else if strings.Compare( string( buf[:n] ), "fin" ) == 0 {
            // debug print
//            fmt.Println( "premature file closing" )
            file.Close()    // close file
            // possibly strike exit...
            os.Exit( 2 )        // quit
        // if conn.Read is not start or fin, it's data
        } else {
            // write info to file
            _, err = writer.Write( buf[:n] )            
            check( err )            // make sure write worked
            err = writer.Flush()    // flush write to file
            check( err )            // make sure write worked
            // debug print
//            fmt.Println( "writing" )
        }
    }
}


var file *os.File
var writer *bufio.Writer


func main() {

    // ip:port
    // ip doesn't exist, implies localhost
    service := ":1202"      // capture ip address and host
    
    // "resolve" ip & host according to TCP rules
    tcpAddr, err := net.ResolveTCPAddr( "tcp", service )
    check( err )    // check error

    // bind and "listen" to ip and port, according to tcp rules
    listener, err := net.ListenTCP( "tcp", tcpAddr )
    check( err )    // check error


    // iterate forever
    // TODO i mean i can totally make this more user friendly
    for {
        // accept a connection that makes its way to bound port
        conn, err := listener.Accept()
        if err != nil {     // if connection fails...
            continue        // don't quit program, not fatal error
        }

        // asynchronous function to handle connection to client
        go handleClient( conn )
    }
}



